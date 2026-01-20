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

type alias Bounds =
    { minX : Float
    , maxX : Float
    , minY : Float
    , maxY : Float
    }


-- PUBLIC API

display : String -> List Instruction -> Svg msg
display color program =
    let
        ( segments, b ) =
            simulate program

        -- On ajoute une marge de sécurité (padding) de 20 unités
        padding = 20
        
        -- Calcul de la largeur et hauteur réelles du tracé
        drawWidth = (b.maxX - b.minX) + (padding * 2)
        drawHeight = (b.maxY - b.minY) + (padding * 2)

        -- Construction de la viewBox : "minX minY largeur hauteur"
        viewBoxStr =
            String.fromFloat (b.minX - padding)
                ++ " "
                ++ String.fromFloat (b.minY - padding)
                ++ " "
                ++ String.fromFloat drawWidth
                ++ " "
                ++ String.fromFloat drawHeight
    in
    svg
        [ width "600"
        , height "600"
        , viewBox viewBoxStr
        , style "background-color: white; border: 1px solid #ddd;"
        ]
        (List.map (segmentToSvg color) segments)


-- SIMULATION DE LA TORTUE

simulate : List Instruction -> ( List Segment, Bounds )
simulate instructions =
    let
        startPos = ( 0, 0 )
        startAngle = 0
        initialBounds = { minX = 0, maxX = 0, minY = 0, maxY = 0 }
    in
    simulateHelp instructions startPos startAngle [] initialBounds


simulateHelp :
    List Instruction
    -> ( Float, Float )
    -> Float
    -> List Segment
    -> Bounds
    -> ( List Segment, Bounds )
simulateHelp instructions ( x, y ) angle acc bounds =
    case instructions of
        [] ->
            ( List.reverse acc, bounds )

        instr :: rest ->
            case instr of
                Forward n ->
                    let
                        rad = degrees angle
                        dx = toFloat n * cos rad
                        dy = toFloat n * sin rad
                        newX = x + dx
                        newY = y + dy
                        
                        newSegment = { x1 = x, y1 = y, x2 = newX, y2 = newY }
                        
                        -- On utilise Basics.min et Basics.max pour lever l'ambiguïté
                        newBounds =
                            { minX = Basics.min bounds.minX newX
                            , maxX = Basics.max bounds.maxX newX
                            , minY = Basics.min bounds.minY newY
                            , maxY = Basics.max bounds.maxY newY
                            }
                    in
                    simulateHelp rest ( newX, newY ) angle (newSegment :: acc) newBounds

                Left n ->
                    simulateHelp rest ( x, y ) (angle - toFloat n) acc bounds

                Right n ->
                    simulateHelp rest ( x, y ) (angle + toFloat n) acc bounds

                Repeat k block ->
                    if k <= 0 then
                        simulateHelp rest ( x, y ) angle acc bounds
                    else
                        simulateHelp (block ++ (Repeat (k - 1) block :: rest)) ( x, y ) angle acc bounds

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
        , strokeLinecap "round"
        ]
        []
