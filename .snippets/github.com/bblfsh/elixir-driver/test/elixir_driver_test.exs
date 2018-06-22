defmodule ElixirDriverTest do
  use ExUnit.Case

  doctest ElixirDriver

  test "map_for_error returns status error and an array of reasons" do
    reason = "I want it to fail"
    result = ElixirDriver.map_for_error reason
    assert_has_key_with_value result, "status", "error"

    assert Map.has_key? result, "errors"
    assert Enum.member? result["errors"], reason
  end

  test "map_for_ast returns status ok and the ast passed" do
    ast = {:my_ast, ["context"], {:nested_ast}}
    result = ElixirDriver.map_for_ast ast
    assert_has_key_with_value result, "status", "ok"
    assert_has_key_with_value result, "ast", ast
  end

  test "driver_info_map contains current version of project" do
    result = ElixirDriver.driver_info_map
    assert Map.has_key? result, "driver"

    driver = result["driver"]
    assert String.contains? driver, "elixir"
    assert String.contains? driver, "0.1.0"
  end

  test "default_map_from_req returns language and version from req" do
    lang = "lang"
    version = "version"
    req = %{"language" => lang, "language_version" => version }
    result = ElixirDriver.default_map_from_req req
    assert_has_key_with_value result, "language", lang
    assert_has_key_with_value result, "language_version", version
  end

  test "default_map_from_req defaults language to 'elixir'" do
    version = "version"
    req = %{ "language_version" => version }
    result = ElixirDriver.default_map_from_req req
    assert_has_key_with_value result, "language", "elixir"
    assert_has_key_with_value result, "language_version", version
  end

  test "default_map_from_req defaults language_version to current elixir version" do
    lang = "lang"
    req = %{ "language" => lang }
    result = ElixirDriver.default_map_from_req req
    assert_has_key_with_value result, "language", lang
    assert_has_key_with_value result, "language_version", System.version()
  end

  test "write_out packs the message and writes it" do
    msg = %{"key" => ["value"]}
    contents = ElixirDriver.pack(msg)

    {:ok, unpacked} = Poison.Parser.parse contents
    assert unpacked == msg
  end

  test "make_fatal changes status to fatal only if it is error" do
    fatal = ElixirDriver.make_fatal %{"status" => "error"}
    assert_has_key_with_value fatal, "status", "fatal"

    no_error = ElixirDriver.make_fatal %{"status" => "ok"}
    assert no_error["status"] == "fatal"

    no_status = ElixirDriver.make_fatal %{"other" => "thing"}
    assert no_status["status"] == "fatal"
  end

  test "build_response returns a valid response for a valid content" do
    %{"status" => "ok", "ast" => ast} = ElixirDriver.build_response "defmodule MyModule do end"
    { definition, _, _ } = ast
    assert definition == :defmodule
  end

  test "build_response returns a valid response for an invalid content" do
    %{"status" => "error", "errors" => errors} = ElixirDriver.build_response "defmodule MyModule  end"
    assert length(errors) > 0
  end

  test "build_response returns a valid response for a working elixir source file" do
    %{"status" => "ok", "ast" => ast} = ElixirDriver.build_response valid_ex_file()
    { definition, _, _ } = ast
    assert definition == :defmodule
  end

  test "handle error returns a packed error" do
    reason = "my reason"
    json = ElixirDriver.handle {:error, reason}
    {:ok, resp} = Poison.Parser.parse json

    assert_response resp, "fatal"
    assert_has_key_with_value resp, "errors", [reason]
  end

  test "handle ok request returns a valid response" do
    req = %{"action" => "ParseAST", "language" => "elixir", "content" => valid_ex_file()}
    json = ElixirDriver.handle {:ok, req}
    {:ok, resp} = Poison.Parser.parse json

    assert_response resp
    assert Map.has_key? resp, "ast"

  end

  test "return_error for estale returns a message about NFS" do
    resp = ElixirDriver.return_error :estale, %{}

    assert_response resp, "error"
    assert String.contains?(hd(resp["errors"]), "NFS")
  end

  test "process with an unknown action returns an error" do
    resp = ElixirDriver.process(%{"action" => "UnknownActionForTest"})

    assert_response resp, "error"
    assert String.contains?(hd(resp["errors"]), "UnknownActionForTest")
  end

  test "process with unknown command errors" do
    resp = ElixirDriver.process("Calimero")

    assert_response resp, "error"
    assert String.contains?(hd(resp["errors"]), "Calimero")
  end

  def assert_response(response, status \\ "ok") when is_map(response) do
    assert_has_key_with_value response, "status", status
    assert_has_key_with_value response, "language", "elixir"
    assert_has_key_with_value response, "language_version", System.version()
    assert_has_key_with_value response, "driver", "elixir #{ElixirDriver.Mixfile.project[:version]}"
  end

  def assert_has_key_with_value map, key, value do
    assert Map.has_key?(map, key), "key #{key} could not be found in map with keys #{map |> Map.keys |> Enum.sort |> Enum.join(", ")} (map: #{inspect map})"
    assert map[key] == value
  end

  def valid_ex_file do
    File.read! "#{File.cwd!}/lib/elixir_driver.ex"
  end
end
