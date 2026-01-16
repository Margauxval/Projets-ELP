module ParserTcTurtle exposing (Instruction(..), read)

import Parser exposing (Parser, (|.), (|=))
import Parser as P


-- TYPES

type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)


-- PUBLIC API

read : String -> Result String (List Instruction)
read source =
    case P.run program source of
        Ok prog ->
            Ok prog

        Err _ ->
            Err "Programme TcTurtle invalide."


-- PARSER

program : Parser (List Instruction)
program =
    P.spaces
        |> P.andThen (\_ -> block)


block : Parser (List Instruction)
block =
    P.sequence
        { start = "["
        , separator = ","
        , end = "]"
        , spaces = P.spaces
        , item = P.lazy (\_ -> instruction)
        , trailing = P.Forbidden
        }


instruction : Parser Instruction
instruction =
    P.oneOf
        [ P.succeed Forward
            |. P.keyword "Forward"
            |. P.spaces
            |= P.int

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
            |= P.int
            |. P.spaces
            |= P.lazy (\_ -> block)
        ]