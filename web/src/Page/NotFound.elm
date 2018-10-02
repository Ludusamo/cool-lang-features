module Page.NotFound exposing (view)

import Browser
import Html exposing (text)



-- View


view : (a -> msg) -> Browser.Document msg
view toMsg =
    { title = "Page Not Found"
    , body =
        [ text "404: Page could not be found" ]
    }
