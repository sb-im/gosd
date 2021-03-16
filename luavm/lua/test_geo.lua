function main(plan)
  print("=== RUN Geography RUN ===")

  -- 114.2247765, 22.6857991
  -- 114.22475167, 22.68580217
  -- = 2.57202994
  local distance = Geo:Distance(114.2247765, 22.6857991, 114.22475167, 22.68580217)
  if math.floor(distance) == 2 then
    print("Distance:", distance)
  else
    error("Distance sum error:", distance)
  end

  local distance = GetDistance(114.2247765, 22.6857991, 114.22475167, 22.68580217)
  if math.floor(distance) == 2 then
    print("Distance:", distance)
  else
    error("Distance sum error:", distance)
  end

  print("=== END Geography END ===")
end
