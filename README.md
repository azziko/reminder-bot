# Telegram reminder bot

Have you ever forgotten to read an article you saved? Well, you are not alone then.
In fact, my telegram is full of different links I leave for "tomorrow", but never actually open them. Since I use telegram a lot, I decided to manage those link within the platform and created the reminder bot. The bot serves for both managing links and reminding of them randomly within 24 hours.

## Features
The working principle is quite straightforward. If you found an interesting article(or pretty much any link starting with http(s)://) and can't or don't want to read it for right now, you simply copy the URL to bot so it can remind you of it later. 

![example](https://user-images.githubusercontent.com/101587881/163728629-67e7a409-7dc6-4d23-b87a-5f76b3822fe1.png)


Managing is also made simple. You can just use one of the available commands that are descriptive enough. 
**Note**: when you call to */rnd* or get a reminder, the retrieved link is immediately deleted from database. 

![commands](https://user-images.githubusercontent.com/101587881/163728644-a515b147-14e7-4517-8e12-ddeaf64d9a0a.png)
## Under the hood
The project is written the way that you can actually use it for any platform(discord, vk and etc.) and any kind of database as long as all corresponding parts implement their interfaces.  
For simplicity sake I implemented polling approach calling getUpdate method every second.  
Was struggling with adding the reminder feature at first, but then came up with an awkward yet effective solution - create a function and execute it concurrently with `time.Sleep()`

To test the project out on your local machine(in case I have not hosted the bot yet) follow these steps:  
1) Pull or fork the repository  
2) cd to the repository in cmd  
3) Run `go build` and then `./remindbot -token <PASS_YOUR_TOKEN_HERE>`


## Plans
Still trying to figure out a better way to implement reminder function, while maintaining the simplicity of polling approach. The current solution I came up with works well enough, but it feels quite inappropriate. That said, it passes an internal race test, so I'm fine with it.  
If you have any suggestions, don't hesitate to drop me a message or contribute.

## License
[MIT](https://choosealicense.com/licenses/mit/)
