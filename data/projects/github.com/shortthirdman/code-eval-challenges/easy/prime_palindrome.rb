require 'prime'
909.step(101, -202) do |i|
  90.step(0, -10) do |j|
    next unless Prime.prime?(i + j)
    puts i + j
    exit 0
  end
end
