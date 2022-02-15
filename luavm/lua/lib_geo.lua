
Geo = {}

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

function Geo:Distance(aLng, aLat, bLng, bLat)
  local range = function(x)
    -- https://pkg.go.dev/math#Acos
    -- Acos(x) = NaN if x < -1 or x > 1
    local min = -1
    local max = 1
    if x > max then
      return max
    elseif x < min then
      return min
    else
      return x
    end
  end

  -- Earth Radius: 6371004
  return 6371004 * math.acos(
    range(
      math.sin(math.rad(aLat)) * math.sin(math.rad(bLat))
      +
      math.cos(math.rad(aLat)) * math.cos(math.rad(bLat))
      *
      math.cos(math.rad(bLng - aLng))
    )
  )
end

function GetDistance(aLng, aLat, bLng, bLat)
  return Geo:Distance(aLng, aLat, bLng, bLat)
end

