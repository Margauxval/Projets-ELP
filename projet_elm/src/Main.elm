
module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Svg exposing (..)
import Svg.Attributes exposing (..)

{- main =
  svg
    [ width "120"
    , height "120"
    , viewBox "0 0 120 120"
    ]
    [ rect
        [ x "10"
        , y "10"
        , width "100"
        , height "100"
        , rx "15"
        , ry "15"
        ]
        []
    , circle
        [ cx "50"
        , cy "50"
        , r "50"
        ]
        []
    ] -}

-- MAIN

main =
    Browser.sandbox { init = init, update = update, view = view }


-- MODEL

type alias Model =
    Int


init : Model
init =
    0


-- UPDATE

type Msg
    = Couleur
    | Avancer
    | Tourner
    | Cercle
    | Triangle
    | Tortue
    | Rectangle


update : Msg -> Model -> Model
update msg model =
    case msg of
        Couleur ->
            model + 1

        Avancer ->
            model + 1

        Tourner ->
            model + 1

        Cercle ->
            model + 1

        Triangle ->
            model + 1

        Tortue ->
            model + 1

        Rectangle ->
            model + 1


-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ button [ onClick Couleur ] [ Html.text "Changer de couleur" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Avancer ] [ Html.text "Avancer" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Tourner ] [ Html.text "Tourner" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Cercle ] [ Html.text "Cercle" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Triangle ] [ Html.text "Triangle" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Tortue ] [ Html.text "Tortue" ]
        , div [] [ Html.text (String.fromInt model) ]

        , button [ onClick Rectangle ] [ Html.text "Rectangle" ]
        , div [] [ Html.text (String.fromInt model) ]
        ]
