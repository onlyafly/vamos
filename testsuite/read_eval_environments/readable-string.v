(list
  (read-string "(1 2 3)")
  (readable-string '(1 2 3))
  (readable-string '(a \b 'c d))
  (read-string (readable-string '(1 2 3))))
