
--[[
6371004*ACOS(
        (
                SIN(RADIANS(C2))*SIN(RADIANS(F2))
                +
                COS(RADIANS(C2))*COS(RADIANS(F2))
                *
                COS(RADIANS(E2-B2))
        )
)
--]]

function GetDistance(aLng, aLat, bLng, bLat)
  -- Earth Radius: 6371004
  return 6371004 * math.acos(
      math.sin(math.rad(aLat)) * math.sin(math.rad(bLat))
      +
      math.cos(math.rad(aLat)) * math.cos(math.rad(bLat))
      *
      math.cos(math.rad(bLng - aLng))
    )
end


