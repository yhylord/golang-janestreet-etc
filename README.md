# 2020 Jane Street ETC Reflection

## 1st stage - 2 hrs

### Strategy

#### Intuition

It would be a good start if we could get anything <ins>easy</ins> and <ins>safe</ins> running at the beginning. Thus we imployed <ins>penny pinching</ins> strategy. 

#### Specifics

There exists one `BOND` with <ins>stable fair value</ins> of <ins>1000</ins>. Thus we place orders of `Dir: "BUY", Price: 999, ...` and `Dir: "SELL", Price: 1001, ...`. After an arbitrary order is filled, we place the order again. Because the bots would buy/sell orders that earliest get to exchange, 1)faster language has some advantage. 2)size doesn't matter that much. 

#### P/L

100~200 per round

### Tech

####  TCP Connection

Jane Street provides written TCP clients for `Python` `C/C++` `Java` `Ocaml` 

#### Programming Language Choice

We chose `Go` as our programming language for following reasons:

​	1) faster than `Python` 

​	2) easier than `C/C++` and `Ocaml`  

​	3) less cumbersome than `Java`

#### Bug

- uppercase and lowercase
  - When we write to exchange intending to place an order, we are required to enter our `Dir` (direction) , which according to the handbook, can be either `"BUY"` or `"SELL"` .
  - Instead we entered `"buy"` and `"sell"` at the beginning and it took us (and TA) 15 mins around to figure it out
- type sensitive 
  - We need to indicate info from exchange as `float64` and then convert it to `int` *a lot* 
- restart program
  - Each round lasts for 5 minutes. The program will very likly be doing nothing once an exchange restarts.
  - In the `connect` function, right after the socket is created (`s = socket.socket(...)`), add this line: `s.settimeout(10)`. This will cause your app to automatically close if the exchange doesn't send any message for 10 seconds.

## 2nd Stage - 3 hrs

### Strategy

#### Intuition

At this stage, we move to address <ins>ADR</ins> to make more profit. ADR has 1 <ins>underlying</ins>(denoted as <ins>U</ins>) security(basically it's an ETF with one underlying component). Since U is more liquid, it might have slower volatility, and because ADR has close connection with U, we could use U to predict ADR. We could just observe U and buy/sell ADR to slightly avoid the pain of conversion fee. 

#### Specifics

- fair value

  - the mean of the highest buying price of U and the lowest selling price of U as our fair value of U. (If the intuition is right, it should also be quite close to the fair value of ADR)

- trades

  - if the highest buying price of U > fair value, sell ADR at (fair value + margin)
  - if the lowest selling price of U < fair value, buy ADR (fair value - margin)

  The limit of ADR is 10 so you might hit the limit pretty fast. Don't panic. Since we only buy ADR when the fair value is higher, i.e. the price will rise in the future, we'll make profit. 

- margin

  - At first, our margin was the midpoint between the higest/lowest price and fair value
  - we adjusted it to 75th percentile and it improved our performance a bit
  - we did not have time to implement but might be useful: analyze the highest/lowest price of **executed** orders to decide margin and maybe take 95 percentile so that you maximize your room to profit in a safe way. (The previous price information was all taken from the book)

#### P/L

300~600 per round

### Tech

We circled a bit to figure out an appropriate way to implement it. The followings are some ideas we had at the first place. They might be helpful if Jane Street change some rules or your situation is more appropriate. If you want to see the solution we ended up, view **solution**. 

We at the beginning were stuck at passing information read from exchange to different file as we cannot have multiple TCP connections at the same time(i.e. we cannot run multiple `ReadFromExchange`  at the same time). Then we thought about implementing an array of channels `[]chan` and send information to those channels so that each strategy could just read information from channels instead of reading from exchange again. We discarded this thought because it was a bit troublesome to implement due to `deadlock` of Go. 

We were concerned that there would be only a few price information in each book as the book did not have accumulated information. We thought about two solutions:

​	1) wait until we have a certain amount of 1)books 2)price information to calculate fair value & make a decision

We take <ins>neither</ins> after the consideration:

​	If we wait a while to make a decision, the information might be outdated.

#### solution

- Two books (no more)

  so that the information is the latest possible and we have just enough to make a buy/sell ADR decision 

  - one for ADR
  - one for U

- Strategy struct 

  - We pass a strategy pointer from the main function to the strategy function and two books are stored in strategy. Basically, all information needed are within the struct so that we don't actually need to store any information in other files. 

## 3rd Stage - 1hr

### Strategy

#### Intuition

ETF is composed by 4 different underlyings at the total amount of 10 (i.e. each underlying has different size). We still did not want to deal with conversion fee neither did we have any time to do so. As a result we decided to employ the same strategy as in 2nd Stage, plus considering the weights of each underlying since they have different size.

#### Specifics

It is basically the same as ADR but with weights. For instance, the weight of U is 1 since there is only one such underlying. Assuming each underlying in ETF has size 2, 2, 2, 4, then the weights could just be exactly the same as the size. If you have time, maybe you could adjust the weight a bit, but based on our experiement, using size as a only factor is sufficient. 

#### P/L

4000 per round

### Tech

We intended to merge the implementations of ADR and ETF as they were quite similiar thus we decided to build a relatively abstract structure. We could consider ETF/ADR as the <ins>composite</ins> with <ins>weighted</ins> <ins>underlying(s)</ins>.

We pass in an array of underlying(s), an array of weight(s), string of composite to the function to create an abstract strategy. As we have an array of books in our struct, we record the book we needed and count the books already obtained, based on the length of underlying(s) and the count of obtained books, we decide whether we would make a trade deicison. (The code would be clearer.)

#### Bug

It might be helpful if you run the strategy in test first and print the fair value and prices. The basket of ETF contains more securities and each might have relatively extrem prices, so the fair value might be no where near the prices of ETF. You could choose to manually set some magic number to make the fair value closer to ETF prices.

- clear operation
  - If you write fair values or books_obtained in the struct like we did, remember to clear them (i.e. reset them as 0) after each trade execution. Otherwise, it would carry to the next deicison making. 



## Collaboration

Things about collaboration that we could have paid more attention about: 

- structure 
  - we spent a fair amount of time figuring out the structure and we could have either left it dirty or had an general idea about it at the beginning.
- coding proficiency
  - we were not that familiar with Go but chose it out of the consideration of speed
- collaboration method
  - git add commit push pull are a lot
  - We used live share in VS code but we were not familiar with VS code thus it didn't provide us with syntax highlight
  - We tended to discuss issues all togethe without anyone simultaneously doing anything and we would all stuck for a relatively long time.





































