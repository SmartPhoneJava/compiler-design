function max() return end
local squares = max({5, 6, 7, 8, 9, 10}, function (v) return v ^ 2 end)

max( t1, { __index = function( t, k ) 
  if t2[ k ] ~= nil then
    return t2[ k ]
  else
    return up1 or up2
  end
end } )