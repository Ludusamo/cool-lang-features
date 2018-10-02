module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, div)
import Http
import Page.NotFound as NotFound
import Url
import Url.Parser exposing ((</>), Parser, int, map, oneOf, parse, s, top)


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = LinkClicked
        , onUrlChange = UrlChanged
        }


init : () -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init _ url key =
    ( Model MainPage [] key, Cmd.none )



-- Model


type alias Model =
    { page : Route
    , features : List Feature
    , key : Nav.Key
    }


type alias Feature =
    { id : String
    , name : String
    , description : String
    }



-- Update


type Msg
    = LoadingFeatures (Result Http.Error (List Feature))
    | LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

        LoadingFeatures (Ok features) ->
            ( { model | features = features }, Cmd.none )

        LoadingFeatures (Err err) ->
            case err of
                Http.BadUrl errorMsg ->
                    Debug.log
                        errorMsg
                        ( model, Cmd.none )

                Http.Timeout ->
                    Debug.log
                        "Timeout"
                        ( model, Cmd.none )

                Http.NetworkError ->
                    Debug.log
                        "NetworkError"
                        ( model, Cmd.none )

                Http.BadStatus status ->
                    Debug.log
                        status.body
                        ( model, Cmd.none )

                Http.BadPayload status res ->
                    Debug.log
                        status
                        ( model, Cmd.none )

        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model
                    , Nav.pushUrl model.key (Url.toString url)
                    )

                Browser.External href ->
                    ( model
                    , Nav.load href
                    )

        UrlChanged url ->
            selectRoute url model



-- Subscriptions


subscriptions : model -> Sub Msg
subscriptions _ =
    Sub.none



-- View


view : Model -> Browser.Document Msg
view model =
    case model.page of
        NotFound ->
            NotFound.view never

        MainPage ->
            NotFound.view never

        AddFeature ->
            NotFound.view never

        EditFeature _ ->
            NotFound.view never



-- Router


type Route
    = MainPage
    | AddFeature
    | EditFeature Int
    | NotFound


routeParser : Parser (Route -> a) a
routeParser =
    oneOf
        [ map MainPage top
        , map AddFeature (s "add")
        , map EditFeature (s "edit" </> int)
        ]


selectRoute : Url.Url -> Model -> ( Model, Cmd Msg )
selectRoute url model =
    case parse routeParser url of
        Just answer ->
            ( { model | page = answer }
            , Cmd.none
            )

        Nothing ->
            ( { model | page = NotFound }
            , Cmd.none
            )



-- Http
