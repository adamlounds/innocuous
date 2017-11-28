# innocuous

To anyone looking, this is just me playing with Go, honest. Should have used a
private repo on bitbucket. Oh well.

my go-fu is clearly rusty, I couldn't get the telnet server to post to the
channel. I tried creating the channel in main() but that didn't work so I tried
putting it in the server and passing it *back* to main - but that didn't work
either.

My original plan was to have a map/reduce thing going on. The telnet server
receives words and splits them, then passes the slice to a worker channel
before replying to the client.  The worker channel then reduces the full array
into a map of word vs count and "letter" (unicode combining characters etc
notwithstanding) vs count, which it can then pass to other channels which are
the only ones allowed to update the global maps.

Next steps:
  * ask someone for help with channels
  * implement the multi-channel updater
  * first pass at generating the top-n list. Sort over all words/letters would
    normally be simplest, but sorting in go is... funky, and it seems a bit
    wasteful
  * write tests - was waiting for list generation before writing http tests
  * look at optimising - don't need to sort, just look at the current top-n
    list and update it.
