function main(plan)
  print("=== RUN Blobs RUN ===")

  print("Blobs:", plan:FileUrl("test_files"))
  print("Blobs:", plan:FileUrl("test_blobs"))

  print("Job Blobs:", plan:JobFileUrl("test_files"))
  print("Job Blobs:", plan:JobFileUrl("test_blobs"))

  sleep("1s")
  print("=== END Blobs END ===")
end
