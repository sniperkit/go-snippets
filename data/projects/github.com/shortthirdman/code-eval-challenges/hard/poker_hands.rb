V = { '2' => 2, '3' => 3, '4' => 4, '5' => 5, '6' => 6, '7' => 7, '8' => 8,
      '9' => 9, 'T' => 10, 'J' => 11, 'Q' => 12, 'K' => 13, 'A' => 14 }.freeze
def card_ranks(cards)
  ranks = cards.map { |x| V[x[0]] }
  ranks.sort!.reverse!
  ranks = [5, 4, 3, 2, 1] if ranks == [14, 5, 4, 3, 2]
  ranks
end

def straight(ranks)
  ranks == [ranks[0], ranks[0] - 1, ranks[0] - 2, ranks[0] - 3, ranks[0] - 4]
end

def flush(hand)
  hand[0][1] == hand[1][1] && hand[0][1] == hand[2][1] &&
    hand[0][1] == hand[3][1] && hand[0][1] == hand[4][1]
end

def kind(n, ranks)
  ranks.each do |i|
    return i if ranks.count(i) == n
  end
  nil
end

def two_pair(ranks)
  h = kind(2, ranks)
  return nil unless h
  ranks.each do |i|
    return [h, i] if ranks.count(i) == 2 && i != h
  end
  nil
end

def poker(h0, h1)
  case h0 <=> h1
  when -1 then 'right'
  when 1 then 'left'
  else 'none'
  end
end

def hand_rank(hand)
  ranks = card_ranks(hand)
  # straight flush
  return [8, ranks.max] if straight(ranks) && flush(hand)
  # 4 of a kind
  return [7, kind(4, ranks), kind(1, ranks)] if kind(4, ranks)
  # full house
  return [6, kind(3, ranks), kind(2, ranks)] if kind(3, ranks) && kind(2, ranks)
  # flush
  return [5, ranks] if flush(hand)
  # straight
  return [4, ranks.max] if straight(ranks)
  # 3 of a kind
  return [3, kind(3, ranks), ranks] if kind(3, ranks)
  # 2 pair
  return [2, two_pair(ranks), kind(1, ranks)] if two_pair(ranks)
  # kind
  return [1, kind(2, ranks), ranks] if kind(2, ranks)
  # high card
  [0, ranks]
end

File.open(ARGV[0]).each_line do |line|
  s = line.split
  puts poker(hand_rank(s[0..4]), hand_rank(s[5..-1]))
end
