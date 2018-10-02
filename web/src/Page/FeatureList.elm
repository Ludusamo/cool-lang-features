module Page.FeatureList exposing (Feature, Model, Msg, featureDecoder, init, update, view)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, a, button, h1, table, td, text, th, thead, tr)
import Html.Attributes exposing (href, style)
import Html.Events exposing (onClick)
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
    { features : List (Maybe Feature) }



-- Init


init : ( Model, Cmd Msg )
init =
    ( Model [], retrieveFeatures )



-- Update


type Msg
    = RetrieveFeatures (Result Http.Error (List (Maybe Feature)))
    | DeletePressed Int
    | DeleteFeature (Result Http.Error ())
    | ModifyPressed Int


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

        DeletePressed id ->
            ( model, deleteFeature id )

        DeleteFeature result ->
            case result of
                Ok _ ->
                    ( model, retrieveFeatures )

                Err _ ->
                    Debug.log
                        "Failed to delete feature"
                        ( model, Cmd.none )

        ModifyPressed id ->
            ( model, Nav.load ("/modify/" ++ String.fromInt id) )



-- View


view : Model -> Skeleton.Details Msg
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
                    , th [] [ text "Control" ]
                    ]
                ]
                (List.map featureRow model.features)
            )
        , a [ href "/add" ] [ text "Add" ]
        ]
    }


featureRow : Maybe Feature -> Html Msg
featureRow maybeFeat =
    case maybeFeat of
        Just feat ->
            tr []
                [ td []
                    [ text (String.fromInt feat.id) ]
                , td []
                    [ text feat.name ]
                , td [] [ text feat.description ]
                , td []
                    [ button
                        [ onClick (DeletePressed feat.id) ]
                        [ text "Delete" ]
                    , button
                        [ onClick (ModifyPressed feat.id) ]
                        [ text "Modify" ]
                    ]
                ]

        Nothing ->
            text ""



-- HTTP


featureDecoder : Decoder Feature
featureDecoder =
    map3 Feature
        (field "id" int)
        (field "name" string)
        (field "description" string)


featureListDecoder : Decoder (List (Maybe Feature))
featureListDecoder =
    Decode.list (Decode.nullable featureDecoder)


retrieveFeatures : Cmd Msg
retrieveFeatures =
    Http.send RetrieveFeatures (Http.get "/api/feature" featureListDecoder)


deleteFeature : Int -> Cmd Msg
deleteFeature id =
    Http.send
        DeleteFeature
        (Http.request
            { body = Http.emptyBody
            , headers = []
            , expect = Http.expectStringResponse (\_ -> Ok ())
            , timeout = Nothing
            , withCredentials = False
            , method = "DELETE"
            , url = "/api/feature/" ++ String.fromInt id
            }
        )
