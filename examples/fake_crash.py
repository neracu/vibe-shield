import sys
sys.stderr.write(
    'Traceback (most recent call last):\n'
    '  File "missing.py", line 5, in <module>\n'
    '    raise ValueError("boom")\n'
    'ValueError: boom\n'
)
raise SystemExit(1)
