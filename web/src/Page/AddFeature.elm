module Page.AddFeature exposing (Model, Msg, init, update, view)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, div, h1, input, text)
import Html.Attributes exposing (href, placeholder, style, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, field, int, map3, string)
import Json.Encode as Encode
import Page.FeatureList exposing (Feature, featureDecoder)
import Skeleton
import Url.Builder exposing (relative)



-- Model


type alias Model =
    { name : String
    , description : String
    , errorMsg : String
    }



-- Init


init : ( Model, Cmd Msg )
init =
    ( Model "" "" "", Cmd.none )



-- Update


type Msg
    = Name String
    | Description String
    | ClickSubmit
    | AddFeature (Result Http.Error Feature)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Name name ->
            ( { model | name = name }, Cmd.none )

        Description description ->
            ( { model | description = description }, Cmd.none )

        ClickSubmit ->
            ( model, addFeature model )

        AddFeature result ->
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



-- View


view : Model -> Skeleton.Details Msg
view model =
    { title = "Add Feature"
    , body =
        [ h1 [] [ text "Add Feature" ]
        , input [ type_ "text", placeholder "Name", value model.name, onInput Name ] []
        , input [ type_ "text", placeholder "Description", value model.description, onInput Description ] []
        , input [ type_ "button", value "Submit", onClick ClickSubmit ] []
        , div [ style "color" "red" ] [ text model.errorMsg ]
        ]
    }



-- HTTP


addFeature : Model -> Cmd Msg
addFeature model =
    Http.send AddFeature
        (Http.post
            "/api/feature"
            (Http.jsonBody
                (Encode.object
                    [ ( "name", Encode.string model.name )
                    , ( "description", Encode.string model.description )
                    ]
                )
            )
            featureDecoder
        )
