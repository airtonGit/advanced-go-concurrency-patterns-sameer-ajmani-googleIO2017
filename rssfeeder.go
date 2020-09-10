package main


type Subscription interface{
  Close() error
}

type sub struct{
  closing chan chan error
}

func (s *sub) Loop(){
  var err error //set when Fetch fails
  var fetchDelay time.Duration // initially 0 (no delay)
  if now := time.Now(); next.After(now) {
    fetchDelay = next.Sub(now)
  }
  startFetch := time.After(fetchDelay) // time.After() return a channel 
  
  var pending []Item // appended by fetch; consumed by send
  var next time.Time // initially January 1, year 0
  
  for{
    var first Item
    var updates chan Item
    if len(pending) > 0 {
      first = pending[0]
      updates = s.updates //habilitar case send
    }
    
    select {
      case errc := <-s.closing:
        errc <- err
        close(s.updates) //tells receiver we're done
        return
      case <-startFetch:
        var fetched []Item
        fetched, next, err = s.fetcher.Fetch()
        if err != nil {
          next = time.Now().Add(10 * time.Second)
          break
        }
        pending = append(pending, fetched...)
      case updates <- first:
        pending = pending[1:]
    }
  }
}


// Pequena explicação
// Envios e recebimentos em canais nil são bloqueantes
// select nunca seleciona um case bloqueante
// Com isso ao setar um canal nil, você pode "desabilitar" um case
// Exemplo de sortear uma moeda
func main() {
  a, b := make(chan string), make(chan string)
  go func() { a <- "a" }()
  go func() { a <- "b" }()
  if rand.Intn(2) == 0 {
    a = nil
    fmt.Println("nil a")
  } else {
    b = nil
    fmt.Println("nil b")
  }
  select {
  case s := <-a:
    fmt.Println("got", s)
  case s := <-b:
    fmt.Println("got", s)
  }
}

func (s *sub) Close(){
  errc := make(chan error)
  s.closing <- errc
  return <- errc
}

