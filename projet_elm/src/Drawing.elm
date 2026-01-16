module Drawing exposing (display)

import Svg exposing (Svg, svg, line)
import Svg.Attributes exposing (..)
import ParserTcTurtle exposing (Instruction(..))


-- TYPES

type alias Segment =
    { x1 : Float
    , y1 : Float
    , x2 : Float
    , y2 : Float
    }


-- PUBLIC API

display : String -> List Instruction -> Svg msg
display color program =
    let
        segments =
            simulate program
    in
    svg
        [ width "600"
        , height "600"
        , viewBox "0 0 600 600"
        ]
        (List.map (segmentToSvg color) segments)


-- SIMULATION DE LA TORTUE

simulate : List Instruction -> List Segment
simulate instructions =
    let
        startPos =
            ( 300, 300 )

        startAngle =
            0
    in
    simulateHelp instructions startPos startAngle [] |> List.reverse


simulateHelp :
    List Instruction
    -> ( Float, Float )
    -> Float
    -> List Segment
    -> List Segment
simulateHelp instructions ( x, y ) angle acc =
    case instructions of
        [] ->
            acc

        instr :: rest ->
            case instr of
                Forward n ->
                    let
                        rad =
                            degrees angle

                        dx =
                            toFloat n * cos rad

                        dy =
                            toFloat n * sin rad

                        newX =
                            x + dx

                        newY =
                            y + dy

                        segment =
                            { x1 = x, y1 = y, x2 = newX, y2 = newY }
                    in
                    simulateHelp rest ( newX, newY ) angle (segment :: acc)

                Left n ->
                    simulateHelp rest ( x, y ) (angle - toFloat n) acc

                Right n ->
                    simulateHelp rest ( x, y ) (angle + toFloat n) acc

                Repeat k block ->
                    let
                        repeated =
                            List.concat (List.repeat k block)
                    in
                    simulateHelp (repeated ++ rest) ( x, y ) angle acc


-- SVG

segmentToSvg : String -> Segment -> Svg msg
segmentToSvg color s =
    line
        [ x1 (String.fromFloat s.x1)
        , y1 (String.fromFloat s.y1)
        , x2 (String.fromFloat s.x2)
        , y2 (String.fromFloat s.y2)
        , stroke color
        , strokeWidth "2"
        ]
        []
