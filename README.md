# Mouse Movement dataset gathering and Generative Adversarial Network training
### The purpose of this application is to gather data to train an ML model to replicate human-like cursor movement behaviour.


TODO:  
1. Data validation in go (at least to enforce proper data structure before inserting into db)
2. Media query? disallow data from touch screen devices - DONE
3. Env variables etc in go to remove sensitive data in repo
4. Pull data from mongo into GAN training set files
5. Train model
6. Add a new endpoint, similar to the circle game that sends the cursor data to the server, which tests the movement behaviour against the GAN
7. Write openCV program that:
    1. Plays an upated version of the circle game, using pyinput or similar 
    2. Uses the GAN model to attempt to move the cursor in a human-like manner

<br>

# Data Gathering:

The front end is a simple frameworkless tool for data gathering, nothing fancy. Its supported by a backend written in GO.
1. Hover the blue circle until it turns green
2. Move the cursor to the red circle, and click anywhere inside it to store the data and reset the circles.
3. Repeat as many times as you wish.

![](example.gif)