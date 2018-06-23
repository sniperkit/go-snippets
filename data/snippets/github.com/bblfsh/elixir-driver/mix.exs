defmodule ElixirDriver.Mixfile do
  use Mix.Project

  def project do
    [app: :elixir_driver,
     version: "0.1.0", # Update in ElixirDriver.driver_info_map/1 as well
     elixir: "~> 1.4",
     build_embedded: Mix.env == :prod,
     start_permanent: Mix.env == :prod,
     deps: deps(),
     preferred_cli_env: [
       "coveralls": :test,
       "coveralls.html": :test,
       "coveralls.json": :test,
     ],
    test_coverage: [tool: ExCoveralls],
    escript: [main_module: ElixirDriverCli,
              name: "native",
              path: "bin/native"]]
  end

  # Configuration for the OTP application
  #
  # Type "mix help compile.app" for more information
  def application do
    # Specify extra applications you'll use from Erlang/Elixir
    [extra_applications: [:logger]]
  end

  # Dependencies can be Hex packages:
  #
  #   {:my_dep, "~> 0.3.0"}
  #
  # Or git/path repositories:
  #
  #   {:my_dep, git: "https://github.com/elixir-lang/my_dep.git", tag: "0.1.0"}
  #
  # Type "mix help deps" for more examples and options
  defp deps do
    [
      {:poison, "~> 3.0.0"},
      {:phst_transform, "~> 1.0.2"},
      {:msgpax, "~> 1.0"},
      {:excoveralls, "~> 0.5.7", only: :test}
    ]
  end
end
