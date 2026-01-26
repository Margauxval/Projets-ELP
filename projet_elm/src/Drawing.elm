module Drawing exposing (display)

import Svg exposing (Svg, svg, line)
import Svg.Attributes exposing (..)
import ParserTcTurtle exposing (Instruction(..))

type alias Segment =
    { x1 : Float
    , y1 : Float
    , x2 : Float
    , y2 : Float
    } 

-- cadrage automatique du dessin
type alias Bounds =
    { minX : Float
    , maxX : Float
    , minY : Float
    , maxY : Float
    }

display : String -> List Instruction -> Svg msg -- transformation liste d'instruction en graphique SVG
display color program =
    let
        -- lancer la simulation pour obtenir les segments et les limites du dessin
        ( segments, b ) =
            simulate program

        padding = 20 -- marge de sécurité autour du dessin

        -- calcul de la taille réelle du contenu dessiné
        drawWidth = (b.maxX - b.minX) + (padding * 2)
        drawHeight = (b.maxY - b.minY) + (padding * 2)

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
        -- transforme chaque segment de calcul en une ligne SVG visible
        (List.map (segmentToSvg color) segments)


-- initialisation
simulate : List Instruction -> ( List Segment, Bounds )
simulate instructions =
    let
        startPos = ( 0, 0 )
        startAngle = 0
        initialBounds = { minX = 0, maxX = 0, minY = 0, maxY = 0 }
    in
    simulateHelp instructions startPos startAngle [] initialBounds


-- traitement récursif des instructions
simulateHelp :
    List Instruction
    -> ( Float, Float ) -- position (x, y) 
    -> Float            -- angle 
    -> List Segment     -- liste des segments accumulés (le dessin)
    -> Bounds           -- limites du dessin
    -> ( List Segment, Bounds )
simulateHelp instructions ( x, y ) angle acc bounds =
    case instructions of
        [] ->
            ( List.reverse acc, bounds )

        instr :: rest ->
            case instr of
                Forward n -> 
                    let
                        -- calcul trigonométrique de la nouvelle position
                        rad = degrees angle
                        dx = toFloat n * cos rad
                        dy = toFloat n * sin rad
                        newX = x + dx
                        newY = y + dy
                        
                        -- création d'un nouveau segment de ligne
                        newSegment = { x1 = x, y1 = y, x2 = newX, y2 = newY }

                        -- mise à jour des limites du cadre (pour le zoom auto)
                        newBounds =
                            { minX = Basics.min bounds.minX newX
                            , maxX = Basics.max bounds.maxX newX
                            , minY = Basics.min bounds.minY newY
                            , maxY = Basics.max bounds.maxY newY
                            }
                    in
                    -- continue avec le reste des instructions
                    simulateHelp rest ( newX, newY ) angle (newSegment :: acc) newBounds

                Left n -> 
                    simulateHelp rest ( x, y ) (angle - toFloat n) acc bounds

                Right n -> 
                    simulateHelp rest ( x, y ) (angle + toFloat n) acc bounds

                Repeat k block -> 
                    if k <= 0 then
                        simulateHelp rest ( x, y ) angle acc bounds
                    else
                        -- ajoute le bloc d'instructions devant le reste, k fois
                        simulateHelp (block ++ (Repeat (k - 1) block :: rest)) ( x, y ) angle acc bounds



-- transforme un objet Segment interne en une balise <line /> HTML/SVG
segmentToSvg : String -> Segment -> Svg msg
segmentToSvg color s =
    line
        [ x1 (String.fromFloat s.x1)
        , y1 (String.fromFloat s.y1)
        , x2 (String.fromFloat s.x2)
        , y2 (String.fromFloat s.y2)
        , stroke color              -- couleur choisie dans l'interface
        , strokeWidth "2"           -- épaisseur du trait
        , strokeLinecap "round"     -- extrémités arrondies 
        ]
        []