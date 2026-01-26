module ParserTcTurtle exposing (Instruction(..), read)

import Parser exposing (Parser, (|.), (|=))
import Parser as P

type Instruction
    = Forward Int                    -- avancer d'une distance X
    | Left Int                       -- tourner à gauche de X degrés
    | Right Int                      -- tourner à droite de X degrés
    | Repeat Int (List Instruction)  -- répéter X fois une liste d'instructions

read : String -> Result String (List Instruction)
read source =
    case P.run program source of
        Ok prog ->
            Ok prog -- renvoie liste instructions

        Err _ ->
            -- Si le texte ne respecte pas les règles 
            Err "Programme TcTurtle invalide."

program : Parser (List Instruction)
program =
    P.spaces -- ignore espaces avant le début
        |> P.andThen (\_ -> block)

block : Parser (List Instruction)
block =
    P.sequence
        { start = "["            
        , separator = ","        
        , end = "]"              
        , spaces = P.spaces      -- autorise les espaces/retours à la ligne n'importe où
        , item = P.lazy (\_ -> instruction) -- ce qu'il y a dedans est une "instruction"
        , trailing = P.Forbidden -- interdit une virgule après la dernière commande
        }

instruction : Parser Instruction
instruction =
    P.oneOf -- le parser essaie chaque option l'une après l'autre
        [ 
          P.succeed Forward
            |. P.keyword "Forward" -- attend le mot exact
            |. P.spaces            -- attend un ou plusieurs espaces
            |= P.int               -- lit un nombre entier

        , P.succeed Left
            |. P.keyword "Left"
            |. P.spaces
            |= P.int

        , P.succeed Right
            |. P.keyword "Right"
            |. P.spaces
            |= P.int

        , P.succeed Repeat
            |. P.keyword "Repeat"
            |. P.spaces
            |= P.int               -- lit le nombre de répétitions
            |. P.spaces
            |= P.lazy (\_ -> block) -- lit ensuite un nouveau bloc (récursivité)
        ]

