module Main exposing (main)

import Browser
import Html exposing (Html, button, div, h2, textarea, text, input)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import ParserTcTurtle exposing (Instruction(..), read)
import Drawing exposing (display)


type alias Model =
    { source : String
    , result : Result String (List Instruction)
    , color : String
    , showDrawing : Bool
    , customColor : String
    }


init : Model
init =
    { source = "[ ]"
    , result = Ok []
    , color = "black"
    , showDrawing = False
    , customColor = ""
    }

type Msg
    = UpdateSource String
    | AddForward
    | AddLeft
    | AddRight
    | AddSquare
    | AddTriangle
    | AddCircle
    | AddHeart
    | AddStar
    | ClearProgram
    | Undo
    | SetColor String
    | Trace
    | UpdateCustomColor String
    | ApplyCustomColor
    | RandomColor

update : Msg -> Model -> Model
update msg model =
    let
        autoRun newModel =
            { newModel | result = read newModel.source }
    in
    case msg of
        UpdateSource txt ->
            { model | source = txt }

        AddForward ->
            { model | source = appendInstruction "Forward 20" model.source }

        AddLeft ->
            { model | source = appendInstruction "Left 15" model.source }

        AddRight ->
            { model | source = appendInstruction "Right 15" model.source }

        AddSquare ->
            { model | source = appendInstruction "Repeat 4 [Forward 80, Left 90]" model.source }

        AddTriangle ->
            { model | source = appendInstruction "Repeat 3 [Forward 100, Left 120]" model.source }

        AddCircle ->
            { model | source = appendInstruction "Repeat 360 [Forward 2, Left 1]" model.source }

        AddHeart ->
            { model | source = appendInstruction "Left 45, Repeat 90 [ Forward 2, Left 2 ], Right 90, Repeat 90 [ Forward 2, Left 2 ], Left 1, Forward 118, Left 90, Forward 118" model.source }

        AddStar ->
            { model | source = appendInstruction "Repeat 5 [Forward 150, Right 144]" model.source }

        ClearProgram ->
            { model | source = "[ ]" }

        Undo ->
            { model | source = undoInstruction model.source }

        SetColor col ->
            { model | color = col }

        Trace ->
            autoRun { model | showDrawing = True }

        UpdateCustomColor txt ->
            { model | customColor = txt }

        ApplyCustomColor ->
            { model | color = model.customColor }

        RandomColor ->
            let
                colors =
                    [ "#FF5733", "#33FF57", "#3357FF", "#FFD700", "#FF69B4", "#12ABB4" ]

                index =
                    modBy (List.length colors) (String.length model.source)

                randomChoice =
                    case List.drop index colors |> List.head of
                        Just c ->
                            c

                        Nothing ->
                            "black"
            in
            { model | color = randomChoice }


appendInstruction : String -> String -> String
appendInstruction instr source =
    case source of
        "[ ]" ->
            "[ " ++ instr ++ " ]"

        _ ->
            let
                inside =
                    String.dropLeft 1 (String.dropRight 1 source)
            in
            "[ " ++ inside ++ ", " ++ instr ++ " ]"


undoInstruction : String -> String
undoInstruction source =
    case source of
        "[ ]" ->
            source

        _ ->
            let
                inside =
                    String.dropLeft 1 (String.dropRight 1 source)

                items =
                    String.split "," inside
                        |> List.map String.trim

                newItems =
                    case List.reverse items of
                        [] ->
                            []

                        _ :: rest ->
                            List.reverse rest
            in
            case newItems of
                [] ->
                    "[ ]"

                _ ->
                    "[ " ++ String.join ", " newItems ++ " ]"


view : Model -> Html Msg
view model =
    div
        [ style "padding" "20px"
        , style "font-family" "'Quicksand', sans-serif"
        , style "background-color" "#d8ecff"
        , style "min-height" "100vh"
        ]
        [ h2 [] [ text "Magic's Draw" ]

        , div
            [ style "display" "flex"
            , style "gap" "40px"
            , style "align-items" "flex-start"
            ]
            [
              div []
                [ h2 [] [ text "Mouvement" ]
                , blueButton AddForward "Avancer 20"
                , blueButton AddLeft "Gauche 15°"
                , blueButton AddRight "Droite 15°"

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Formes" ]
                    , blueButton AddSquare "Carré"
                    , blueButton AddTriangle "Triangle"
                    , blueButton AddCircle "Cercle"
                    , blueButton AddHeart "Cœur"
                    , blueButton AddStar "Étoile"
                    ]

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Couleurs" ]
                    , blueButton (SetColor "black") "Noir"
                    , blueButton (SetColor "red") "Rouge"
                    , blueButton (SetColor "blue") "Bleu"
                    , blueButton (SetColor "green") "Vert"

                    , div [ style "margin-top" "15px" ]
                        [ input
                            [ value model.customColor
                            , onInput UpdateCustomColor
                            , placeholder "Exemple : #12ABB4"
                            , style "color" "gray"
                            , style "font-size" "14px"
                            , style "width" "150px"
                            , style "height" "30px"
                            ]
                            []
                        , blueButton ApplyCustomColor "Valider mon choix"
                        , blueButton RandomColor "Random"
                        ]
                    ]

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Contrôle" ]
                    , blueButton Undo "Annuler"
                    , blueButton ClearProgram "Effacer"
                    , blueButton Trace "Tracer"
                    ]
                ]

            , div []
                [ textarea
                    [ value model.source
                    , onInput UpdateSource
                    , rows 12
                    , cols 55
                    , style "font-size" "16px"
                    ]
                    []

                , if model.showDrawing then
                    div
                        [ style "margin-top" "20px"
                        , style "border" "1px solid #ccc"
                        , style "padding" "10px"
                        , style "display" "inline-block"
                        , style "background" "white"
                        ]
                        [ viewResult model ]
                  else
                    text ""
                ]
            ]
        ]


blueButton : Msg -> String -> Html Msg
blueButton msg label =
    button
        [ onClick msg
        , style "background-color" "#007BFF"
        , style "color" "white"
        , style "border" "none"
        , style "padding" "8px 12px"
        , style "margin" "5px"
        , style "border-radius" "4px"
        , style "cursor" "pointer"
        ]
        [ text label ]


viewResult : Model -> Html msg
viewResult model =
    case model.result of
        Err msg ->
            div [ style "color" "red" ] [ text msg ]

        Ok program ->
            display model.color program


main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }

