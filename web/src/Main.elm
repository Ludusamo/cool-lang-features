module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, div)
import Http
import Page.AddFeature as AddFeature
import Page.FeatureList as FeatureList
import Page.ModifyFeature as ModifyFeature
import Page.NotFound as NotFound
import Skeleton
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
    selectRoute url { page = NotFound, key = key }



-- Model


type alias Model =
    { page : Page
    , key : Nav.Key
    }



-- Update


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | FeatureListMsg FeatureList.Msg
    | AddFeatureMsg AddFeature.Msg
    | ModifyFeatureMsg ModifyFeature.Msg
    | NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

        FeatureListMsg message ->
            case model.page of
                FeatureListPage pageModel ->
                    featureList model (FeatureList.update message pageModel)

                _ ->
                    ( model, Cmd.none )

        AddFeatureMsg message ->
            case model.page of
                AddFeaturePage pageModel ->
                    addFeature model (AddFeature.update message pageModel)

                _ ->
                    ( model, Cmd.none )

        ModifyFeatureMsg message ->
            case model.page of
                ModifyFeaturePage pageModel ->
                    modifyFeature model (ModifyFeature.update message pageModel)

                _ ->
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

        FeatureListPage featurelistModel ->
            Skeleton.view FeatureListMsg (FeatureList.view featurelistModel)

        AddFeaturePage addFeatureModel ->
            Skeleton.view AddFeatureMsg (AddFeature.view addFeatureModel)

        ModifyFeaturePage modifyFeatureModel ->
            Skeleton.view ModifyFeatureMsg (ModifyFeature.view modifyFeatureModel)



-- Router


type Page
    = FeatureListPage FeatureList.Model
    | AddFeaturePage AddFeature.Model
    | ModifyFeaturePage ModifyFeature.Model
    | NotFound


routeParser : Model -> Parser (( Model, Cmd Msg ) -> a) a
routeParser model =
    oneOf
        [ route top (featureList model FeatureList.init)
        , route (s "add") (addFeature model AddFeature.init)
        , route (s "modify" </> int)
            (\id -> modifyFeature model (ModifyFeature.init id))
        ]


featureList : Model -> ( FeatureList.Model, Cmd FeatureList.Msg ) -> ( Model, Cmd Msg )
featureList model ( features, cmds ) =
    ( { model | page = FeatureListPage features }
    , Cmd.map FeatureListMsg cmds
    )


addFeature : Model -> ( AddFeature.Model, Cmd AddFeature.Msg ) -> ( Model, Cmd Msg )
addFeature model ( addModel, cmds ) =
    ( { model | page = AddFeaturePage addModel }
    , Cmd.map AddFeatureMsg cmds
    )


modifyFeature : Model -> ( ModifyFeature.Model, Cmd ModifyFeature.Msg ) -> ( Model, Cmd Msg )
modifyFeature model ( modifyModel, cmds ) =
    ( { model | page = ModifyFeaturePage modifyModel }
    , Cmd.map ModifyFeatureMsg cmds
    )


selectRoute : Url.Url -> Model -> ( Model, Cmd Msg )
selectRoute url model =
    case parse (routeParser model) url of
        Just answer ->
            answer

        Nothing ->
            ( { model | page = NotFound }
            , Cmd.none
            )


route : Parser a b -> a -> Parser (b -> c) c
route parser handler =
    map handler parser
