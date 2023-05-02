module FrontEnd exposing (..)

-- MAIN

import Browser
import Html.Attributes exposing (placeholder, style, value)
import Html exposing (..)
import Html.Events exposing (onClick, onInput)
import Http exposing (Error)
import Json.Decode as Decode exposing (Decoder, field, int, map2, string)
import Json.Encode as Encode
import Html exposing (..)
import Json.Encode

main =
  Browser.element
  { init = init
  , update = update
  , subscriptions = subscriptions
  , view = view
  }

-- SUBSCRIPTIONS nonelol
subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none

-- UPDATE

-- MODEL
type alias ToDo =
  { id : Int
  , todo : String
  }

type alias Model =
  { items : List ToDo
  , currentInput : String
  , highestId : Int
  }

init : () -> (Model, Cmd Msg)
init _ =
  ((Model [] "" 0), getToDosHttp)

type Msg
  = InputChange String
  | AddItem
  | RemoveItem ToDo
  -- http calls that are used by above
  | SetAllToDos (Result Http.Error (List ToDo))
  | DeleteToDo (Result Http.Error ToDo)
  | CreateToDo (Result Http.Error ToDo)

-- HTTP functions
getToDosHttp : Cmd Msg
getToDosHttp =
  Http.get
    { url = "http://localhost:9000/v1/todo_service/todo"
    , expect = Http.expectJson SetAllToDos mutlipleToDoDecoder
    }

removeTodoHTTP : Int -> Cmd Msg
removeTodoHTTP id =
  Http.request
    { method = "DELETE"
    , headers = []
    , url = "http://localhost:9000/v1/todo_service/todo/" ++ (String.fromInt id)
    , body = Http.emptyBody
    , expect = Http.expectJson DeleteToDo todoDecoder
    , timeout = Nothing
    , tracker = Nothing
    }

createTodoHTTP: String -> Cmd Msg
createTodoHTTP todoItem =
  let
    body =
      Encode.object[("todo", Encode.string todoItem)]
      |> Http.jsonBody
  in
  Http.request
    { method = "POST"
    , headers = []
    , url = "http://localhost:9000/v1/todo_service/todo"
    , body = body
    , expect = Http.expectJson CreateToDo todoDecoder
    , timeout = Nothing
    , tracker = Nothing
    }

mutlipleToDoDecoder: Decode.Decoder (List ToDo)
mutlipleToDoDecoder =
  Decode.field "todos" (Decode.list (todoDecoder))

todoDecoder : Decoder ToDo
todoDecoder =
  map2 ToDo
    (field "id" int)
    (field "todo" string)


update: Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    InputChange newInput ->
      ({ model | currentInput = newInput }, Cmd.none)

    AddItem ->
      ({ model | items = List.append model.items [ToDo (model.highestId + 1)  model.currentInput]
                 , currentInput = ""
                 , highestId = model.highestId + 1}
                 , createTodoHTTP model.currentInput )

    RemoveItem item ->
      ({ model | items = List.filter (\x -> (x.id /= item.id)) model.items}, removeTodoHTTP item.id)

     -- like on mounted or something ig, runs on init
    SetAllToDos result ->
      case result of
        Ok todos ->
          ({ model | items = List.map (\x -> (ToDo x.id x.todo)) todos
                   , currentInput = ""
                   , highestId = Maybe.withDefault 0 (List.maximum (List.map (\x -> (x.id)) todos))
                   }
          , Cmd.none)

        Err error ->
          ({ model | items = [ToDo -1 (Debug.toString error) ] }, Cmd.none)

    DeleteToDo result ->
      case result of
        Ok item ->
          ({ model | items = List.filter (\x -> (x.id /= item.id)) model.items}, Cmd.none)

        Err error ->
          ({ model | items = [ToDo -1 (Debug.toString error) ] }, Cmd.none)

    CreateToDo result ->
      case result of
        Ok _ ->
          (model , Cmd.none)

        Err error ->
          ({ model | items = [ToDo -1 (Debug.toString error) ] }, Cmd.none)


view : Model -> Html Msg
view model =
  div [ style "margin" "auto"
      , style "width" "max-content"
      , style "padding" "5%"
      , style "font-family" "monospace"
      ]
         [ input [ placeholder "Enter ToDo..."
                , value model.currentInput
                , onInput InputChange]
                []

         , button [ onClick AddItem
                  , style "background-color" "#48c78e"
                  , style "border-color" "#48c78e"
                  , style "color" "white"
                  , style "cursor" "pointer"
                  ]
                  [text "Add ToDo"]

         , h2 [ style "text-decoration" "underline", style "margin-bottom" "5px"] [text "Items"]

         , ul [ style "list-style-type" "none"
              , style "padding" "0 0 0 0"]
              (List.map (\x -> li [ style "margin-top" "10px"] [ button
                                              [ onClick (RemoveItem x)
                                              , style "cursor" "pointer"
                                              , style "border-radius" "5px"
                                              , style "background-color" "#f14668"
                                              , style "color" "white"]
                                              [text "X"]

                                     , div [ style "display" "inline"
                                           , style "padding" "0 10px 0"
                                           ]
                                           [text x.todo]

                ] ) model.items)
  ]
