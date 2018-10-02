module Skeleton exposing (Details, view)

import Browser
import Html exposing (Html, div)


type alias Details msg =
    { title : String
    , body : List (Html msg)
    }


view : (a -> msg) -> Details a -> Browser.Document msg
view toMsg details =
    { title = details.title
    , body = [ Html.map toMsg <| div [] details.body ]
    }
