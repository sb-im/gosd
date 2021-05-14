
_SD_printResult = ""
_SD_raw_print = print
function print (...)
    _SD_raw_print(arg)
    _SD_printResult = _SD_printResult .. os.date("%Y/%m/%d %H:%M:%S") .. " "
    for i,v in ipairs(arg) do
        _SD_printResult = _SD_printResult .. tostring(v) .. "\t"
    end
    _SD_printResult = _SD_printResult .. "\n"
end

