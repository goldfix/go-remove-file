# go-remove-file
Possible substitute at 'rm' command. Similar functions plus recovery and trash functions.

<pre><code>
Usage of grm:
  -d        move files into recycle
  -e        empty recycle folder
  -f        ignore nonexistent files, never prompt
  -i        prompt before every files
  -ls       list files into recycle
  -r        remove contents recursively
  -u        recover list files from recycle
  -v        explain what is being done
  -version  output version information and exit
  -x        force delete files
</code></pre>

Thanks to:
* https://github.com/satori/go.uuid