function main(plan)
  print("=== RUN Extra RUN ===")
  plan:GetExtra("test_key")
  plan:SetExtra("test_key", "test_value")
  sleep("1s")
  print("=== END Extra END ===")
end
