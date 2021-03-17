function main(plan)
  print("=== RUN Notification RUN ===")
  print(os.time())
  plan:Notification("notification")
  sleep("1s")
  plan:Notification("notification", 3)
  sleep("1s")
  plan:Notification("notification", "1")
  sleep("1s")
  plan:Notification(233, "1")

  sleep("1s")
  print("=== END Notification END ===")
end
