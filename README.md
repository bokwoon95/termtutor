# termtutor
NUS Hack&amp;Roll 2019

### How to run
Simply run `./termtutor`. 'data.json' must be in the same directory because that's where `termtutor` sources its data from.
All other files are irrelevant.

### How to use

When it starts up, you are put into 'lesson' mode, with questions on the left and answers on the right.
You can switch to 'lesson' mode anytime with `Ctrl+Q`.
To view the answer to the current question and as well as to show the next question, press Enter.
There are only 10 questions I have come up with so far so they just loop around.
These are pretty rushed questions, merely a proof of concept.
Also note that the alignments starts getting screwy once you go past 10 questions, but I like to spam Enter to see the text in both panes fly by.

To access the Javascript REPL in 'interpreter' mode, press `Ctrl+W`. 
Note that multi-line inputs are not accepted yet because pressing Enter simply sends the current text into the REPL (so you have to write your functions in one line..)
On a macbook you may also experience crashes when using 'backspace' because backspace on a mac is slightly different from the backspace that a terminal recognizes. 
Use `Ctrl+H` instead (synonymous to backspace) if you crash too much.

`Console.log` is borked. It prints the string in the editing textbox rather than in the REPL window, and that text cannot be cleared by `gocui`'s Clear() function.
To get rid of it from your screen, press `Ctrl+O` (which fills up the editing textbox with 'O's, overwriting the `console.log`-ged text)
followed by `Ctrl+U` (which clears the printed 'O's) as a hacky workaround.

To exit the application, press either Ctrl+C or Ctrl+D.
