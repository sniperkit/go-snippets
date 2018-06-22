defmodule BblfshDriver do
  @moduledoc """
  Parses the content and produces the right formatted output
  """

  @doc """
  """
  def parse data do
    data
    |> Code.string_to_quoted
    |> build_response
  end

  def build_response {:ok, quoted_form} do
    {:ok, quoted_form}
  end

  def build_response {:error, {{line, posinit, posend}, error, token}} do
    {:error, "#{token} at #{line} #{posinit} - #{posend} produced error #{error}"}
  end
end
