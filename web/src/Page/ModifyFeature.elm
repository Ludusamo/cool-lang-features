module Page.ModifyFeature exposing (Model, Msg, init, update, view)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, div, h1, input, text)
import Html.Attributes exposing (disabled, href, placeholder, style, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, field, int, map3, string)
import Json.Encode as Encode
import Page.FeatureList exposing (Feature, featureDecoder)
import Skeleton
import Url.Builder exposing (relative)



-- Model


type alias Model =
    { id : Int
    , name : String
    , description : String
    , errorMsg : String
    }



-- Init


init : Int -> ( Model, Cmd Msg )
init id =
    ( Model id "" "" "", retrieveFeature id )



-- Update


type Msg
    = Name String
    | Description String
    | ClickSubmit
    | ModifyFeature (Result Http.Error Feature)
    | ReceiveFeature (Result Http.Error Feature)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Name name ->
            ( { model | name = name }, Cmd.none )

        Description description ->
            ( { model | description = description }, Cmd.none )

        ClickSubmit ->
            ( model, modifyFeature model )

        ModifyFeature result ->
            case result of
                Ok _ ->
                    ( model, Nav.load (relative [ "/" ] []) )

                Err err ->
                    case err of
                        Http.BadStatus res ->
                            ( { model | errorMsg = res.status.message }, Cmd.none )

                        Http.BadPayload res _ ->
                            ( { model | errorMsg = res }, Cmd.none )

                        _ ->
                            ( { model | errorMsg = "Failed to add" }, Cmd.none )

        ReceiveFeature result ->
            case result of
                Ok feat ->
                    ( { model
                        | name = feat.name
                        , description = feat.description
                      }
                    , Cmd.none
                    )

                _ ->
                    ( { model | errorMsg = "Failed to get feature" }, Cmd.none )



-- View


view : Model -> Skeleton.Details Msg
view model =
    { title = "Edit Feature"
    , body =
        [ h1 [] [ text "Modify Feature" ]
        , input [ type_ "text", placeholder "id", value (String.fromInt model.id), disabled True ] []
        , input [ type_ "text", placeholder "Name", value model.name, onInput Name ] []
        , input [ type_ "text", placeholder "Description", value model.description, onInput Description ] []
        , input [ type_ "button", value "Submit", onClick ClickSubmit ] []
        , div [ style "color" "red" ] [ text model.errorMsg ]
        ]
    }



-- HTTP


retrieveFeature : Int -> Cmd Msg
retrieveFeature id =
    Http.send ReceiveFeature (Http.get ("/api/feature/" ++ String.fromInt id) featureDecoder)


modifyFeature : Model -> Cmd Msg
modifyFeature model =
    Http.send
        ModifyFeature
        (Http.request
            { body =
                Http.jsonBody
                    (Encode.object
                        [ ( "id", Encode.int model.id )
                        , ( "name", Encode.string model.name )
                        , ( "description", Encode.string model.description )
                        ]
                    )
            , headers = []
            , expect = Http.expectJson featureDecoder
            , timeout = Nothing
            , withCredentials = False
            , method = "PATCH"
            , url = "/api/feature/" ++ String.fromInt model.id
            }
        )
