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

  -- 118.7476783, 31.9414349
  -- 118.7476783, 31.9414349
  -- = 0
  --
  -- Maybe == NaN
  -- math.acos(0.2798961794852986 + 0.7201038205147016)
  -- https://pkg.go.dev/math#Acos
  -- Acos(x) = NaN if x < -1 or x > 1
  local distance_nan = GetDistance(118.7476783, 31.9414349, 118.7476783, 31.9414349)
  if math.floor(distance_nan) == 0 then
    print("Distance:", distance_nan)
  else
    error("Distance sum error:", distance_nan)
  end

  print("=== END Geography END ===")
end
