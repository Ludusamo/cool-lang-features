module Page.FeatureList exposing (Model, Msg, init, view)

import Browser
import Html exposing (Html, h1, table, td, text, th, thead, tr)



-- Model


type alias Feature =
    { id : String
    , name : String
    , description : String
    }


type alias Model =
    { features : List Feature }



-- Init


init : ( Model, Cmd Msg )
init =
    ( Model [], Cmd.none )



-- Update


type Msg
    = RetrievingFeatures



-- View


view : (a -> msg) -> Model -> Browser.Document msg
view toMsg model =
    { title = "Feature List"
    , body =
        [ h1 [] [ text "Cool Programming Language Features" ]
        , table []
            (List.append
                [ thead []
                    [ th [] [ text "ID" ]
                    , th [] [ text "Name" ]
                    , th [] [ text "Description" ]
                    ]
                ]
                (List.map featureRow model.features)
            )
        ]
    }


featureRow : Feature -> Html msg
featureRow feat =
    tr []
        [ td []
            [ text feat.id
            , text feat.name
            , text feat.description
            ]
        ]
