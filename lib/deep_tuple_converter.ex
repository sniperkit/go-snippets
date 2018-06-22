defmodule DeepTupleConverter do

  def convert tuple do
    PhStTransform.transform(tuple, %{Tuple => fn(t) -> process_tuple(t) end,
                                     Atom => fn(t) -> Atom.to_string(t) end,
                                     Keyword => fn(l) -> process_keyword(l) end
                                   })
  end

  def process_tuple tuple do
    tuple
    |> Tuple.to_list
  end

  def process_keyword keywords do
    keywords
    |> Enum.map(fn {key, value} -> [key, value] end)
  end
end
