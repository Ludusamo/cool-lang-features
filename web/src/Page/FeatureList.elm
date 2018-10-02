module Page.FeatureList exposing (Model, Msg, init, update, view)

import Browser
import Html exposing (Html, h1, table, td, text, th, thead, tr)
import Html.Attributes exposing (style)
import Http
import Json.Decode as Decode exposing (Decoder, field, int, map3, string)
import Skeleton



-- Model


type alias Feature =
    { id : Int
    , name : String
    , description : String
    }


type alias Model =
    { features : List Feature }



-- Init


init : ( Model, Cmd Msg )
init =
    ( Model [], retrieveFeatures )



-- Update


type Msg
    = RetrieveFeatures (Result Http.Error (List Feature))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        RetrieveFeatures result ->
            case result of
                Ok features ->
                    ( { model | features = features }, Cmd.none )

                Err _ ->
                    Debug.log
                        "Failed to retrieve features"
                        ( model, Cmd.none )



-- View


view : Model -> Skeleton.Details msg
view model =
    { title = "Feature List"
    , body =
        [ h1 [] [ text "Cool Programming Language Features" ]
        , table [ style "border" "1px solid black" ]
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
            [ text (String.fromInt feat.id) ]
        , td []
            [ text feat.name ]
        , td [] [ text feat.description ]
        ]



-- HTTP


featureDecoder : Decoder Feature
featureDecoder =
    map3 Feature
        (field "id" int)
        (field "name" string)
        (field "description" string)


featureListDecoder : Decoder (List Feature)
featureListDecoder =
    Decode.list featureDecoder


retrieveFeatures : Cmd Msg
retrieveFeatures =
    Http.send RetrieveFeatures (Http.get "/api/feature" featureListDecoder)
