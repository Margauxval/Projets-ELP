module Main exposing (main)

import Browser
import Html exposing (Html, button, div, h2, textarea, text)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import ParserTcTurtle exposing (Instruction(..), read)
import Drawing exposing (display)


-- MODEL

type alias Model =
    { source : String
    , result : Result String (List Instruction)
    , color : String
    }


init : Model
init =
    { source = "[ ]"
    , result = Ok []
    , color = "black"
    }


-- MESSAGES

type Msg
    = UpdateSource String
    | AddForward
    | AddLeft
    | AddRight
    | AddSquare
    | AddTriangle
    | AddCircle
    | AddTurtle
    | ClearProgram
    | Undo
    | SetColor String


-- UPDATE

update : Msg -> Model -> Model
update msg model =
    let
        autoRun newModel =
            { newModel | result = read newModel.source }
    in
    case msg of
        UpdateSource txt ->
            autoRun { model | source = txt }

        AddForward ->
            autoRun (appendInstruction "Forward 20" model)

        AddLeft ->
            autoRun (appendInstruction "Left 15" model)

        AddRight ->
            autoRun (appendInstruction "Right 15" model)

        AddSquare ->
            autoRun (appendInstruction "Repeat 4 [Forward 80, Left 90]" model)

        AddTriangle ->
            autoRun (appendInstruction "Repeat 3 [Forward 100, Left 120]" model)

        AddCircle ->
            autoRun (appendInstruction "Repeat 360 [Forward 2, Left 1]" model)

        AddTurtle ->
            autoRun
                (appendInstruction
                    """
                    Repeat 1 [
                        Repeat 6 [Forward 60, Left 60],
                        Left 90, Forward 20, Left 180, Forward 20, Left 180,
                        Right 150, Forward 20, Left 180, Forward 20, Left 180,
                        Left 60, Forward 20, Left 180, Forward 20, Left 180,
                        Right 150, Forward 20, Left 180, Forward 20, Left 180,
                        Left 90, Forward 10, Left 180, Forward 10, Left 180
                    ]
                    """
                    model
                )

        ClearProgram ->
            autoRun { model | source = "[ ]" }

        Undo ->
            autoRun (undoInstruction model)

        SetColor col ->
            autoRun { model | color = col }


appendInstruction : String -> Model -> Model
appendInstruction instr model =
    let
        newSource =
            case model.source of
                "[ ]" ->
                    "[ " ++ instr ++ " ]"

                _ ->
                    let
                        inside =
                            String.dropLeft 1 (String.dropRight 1 model.source)
                    in
                    "[ " ++ inside ++ ", " ++ instr ++ " ]"
    in
    { model | source = newSource }


undoInstruction : Model -> Model
undoInstruction model =
    case model.source of
        "[ ]" ->
            model

        _ ->
            let
                inside =
                    String.dropLeft 1 (String.dropRight 1 model.source)

                items =
                    String.split "," inside
                        |> List.map String.trim

                newItems =
                    case List.reverse items of
                        [] ->
                            []

                        _ :: rest ->
                            List.reverse rest

                newSource =
                    case newItems of
                        [] ->
                            "[ ]"

                        _ ->
                            "[ " ++ String.join ", " newItems ++ " ]"
            in
            { model | source = newSource }


-- STYLE POUR LES BOUTONS

buttonStyle : List (Html.Attribute Msg)
buttonStyle =
    [ style "background-color" "#4a90e2"
    , style "color" "white"
    , style "border" "none"
    , style "padding" "10px 20px"
    , style "margin" "5px"
    , style "border-radius" "6px"
    , style "cursor" "pointer"
    , style "font-size" "16px"
    , style "font-family" "'Quicksand', sans-serif"
    , style "box-shadow" "0 4px 6px rgba(0,0,0,0.1)"
    , style "transition" "background-color 0.3s ease"
    ]

styledButton : Msg -> String -> Html Msg
styledButton msg label =
    button (buttonStyle ++ [ onClick msg ]) [ text label ]


-- VIEW

view : Model -> Html Msg
view model =
    div
        [ style "padding" "20px"
        , style "font-family" "'Quicksand', sans-serif"
        , style "background-color" "#d8ecff"
        , style "min-height" "100vh"
        ]
        [ h2 [] [ text "TcTurtle – Interface améliorée" ]

        , div
            [ style "display" "flex"
            , style "gap" "40px"
            , style "align-items" "flex-start"
            ]
            [ -- Colonne gauche : boutons
              div []
                [ h2 [] [ text "Mouvement" ]
                , styledButton AddForward "Forward 20"
                , styledButton AddLeft "Left 15°"
                , styledButton AddRight "Right 15°"

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Formes" ]
                    , styledButton AddSquare "Carré"
                    , styledButton AddTriangle "Triangle"
                    , styledButton AddCircle "Cercle"
                    , styledButton AddTurtle "Tortue"
                    ]

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Couleurs" ]
                    , styledButton (SetColor "black") "Noir"
                    , styledButton (SetColor "red") "Rouge"
                    , styledButton (SetColor "blue") "Bleu"
                    , styledButton (SetColor "green") "Vert"
                    ]

                , div [ style "margin" "20px 0" ]
                    [ h2 [] [ text "Contrôle" ]
                    , styledButton Undo "Undo"
                    , styledButton ClearProgram "Clear"
                    ]
                ]

              -- Colonne droite : zone de texte + dessin
            , div []
                [ textarea
                    [ value model.source
                    , onInput UpdateSource
                    , rows 10
                    , cols 50
                    , style "font-size" "16px"
                    ]
                    []

                , div
                    [ style "margin-top" "20px"
                    , style "border" "1px solid #ccc"
                    , style "padding" "10px"
                    , style "display" "inline-block"
                    , style "background" "white"
                    ]
                    [ viewResult model ]
                ]
            ]
        ]


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
