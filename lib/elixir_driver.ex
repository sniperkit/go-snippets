defmodule ElixirDriver do
  @moduledoc """
  Documentation for ElixirDriver.
  """

  @doc """
  """
  def loop do
    IO.read(:stdio, :line)
    |> Poison.Parser.parse
    |> handle
    |> write
    loop()
  end

  def write item do
    IO.puts :stdio, item
  end

  def return_error reason, req \\ %{}
  def return_error :estale, req do
    return_error "reading from NFS volume", req
  end

  def return_error reason, req do
    map_for_error(reason)
    |> Map.merge(map_for_request(req))
  end

  def handle {:error, reason} do
    return_error(reason)
    |> Map.merge(map_for_request())
    |> make_fatal
    |> pack
  end

  def handle {:ok, req} do
    req
    |> process
    |> build_response
    |> Map.merge(map_for_request())
    |> pack
  end

  def process %{"action" => "ParseAST", "language" => "elixir", "content" => content} do
    content
  end

  def process %{"action" => action} do
    return_error "did not understood action #{action}"
  end

  def process anything do
    return_error "could not understand #{inspect anything}"
  end

  def build_response content do
    case BblfshDriver.parse content do
      {:ok, ast} -> map_for_ast(ast)
      {:error, error} -> map_for_error(error)
    end
  end

  def pack msg do
    case Poison.encode(DeepTupleConverter.convert(msg)) do
      {:ok, packed} -> packed
      {:error, reason} -> pack(return_error reason)
    end
  end

  def map_for_error error do
    %{ "status" => "error", "errors" => [ error ] }
  end

  def map_for_ast ast do
    %{ "status" => "ok", "ast" => ast }
  end

  def map_for_request req \\ %{} do
    Map.merge driver_info_map(), default_map_from_req(req)
  end

  def default_map_from_req req do
    %{
      "language" => lang(req),
      "language_version" => lang_version(req)
    }
  end

  def make_fatal resp do
    Map.put resp, "status", "fatal"
  end

  def driver_info_map do
    %{ "driver" => "elixir 0.1.0" }
  end

  def lang req do
    fetch_with_default(req, "language", "elixir")
  end

  def lang_version req do
    fetch_with_default(req, "language_version", System.version())
  end

  def fetch_with_default map, key, default do
    case Map.fetch(map, key) do
      {:ok, value} -> value
      :error -> default
    end
  end
end
