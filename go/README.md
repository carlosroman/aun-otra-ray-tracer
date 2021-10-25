# Aun Otra Ray Tracer - Golang

## ⚡️ Quick start

First of all, [download](https://golang.org/dl/) and install **Go**. Version `1.16` or higher is required.

Next clone this project and move into the go directory:

```bash
git clone git@github.com:carlosroman/aun-otra-ray-tracer.git
cd aun-otra-ray-tracer/go
```

Then run the following Make command:

```bash
make build/example
```

This will build an example application that renders a couple of simple scenes.
You can find the application in the target directory.
To run it, just run:

```bash
./target/example
```

## 📺 Example application

The example application shows off the ray tracer by rendering a couple of simple scenes.
Running “example” will not run a render but give you the following subcommands you can run:

* cubes
* teapot

### 🎲 Cubes

This command will render some simple geometric shapes to show off the render engine.
By default, it generates an image to a [ppm file](https://en.wikipedia.org/wiki/Netpbm) (`cubes.ppm`).
To output a jpeg image, run the following:

```bash
example cubes -j
```

The file generated will be called `cubes.jpg`. 

This example is one of the quickest to render at the moment.
The default image size is an image that is `640` x `360`.
To render a larger image, just change the width by using the `--width` flag:

```bash
example cubes --width 1080
```

For more options on this subcommand, just run:

```bash
example cubes --help
```

### 🫖 Teapot

This command will render the [Utah teapot](https://en.wikipedia.org/wiki/Utah_teapot) in a simple reflective scene.
The image takes around 30 min to render on an Apple M1. The output will be a file called `teapot.ppm`.
If you want to render a low res version of the teapot, you can use the flag “--low-res”:

```bash
example teapot --low-res
```

The low res teapot takes around one minute to render.
You can speed things up by lowering the sample rate, using the flag `--samples`.
The flag will reduce the number of rays cast per pixel.
You can also use the `--width` flag to render a smaller image as well:

```bash
example teapot --low-res --width 240 --samples 1
```

For more options on this subcommand, just run:

```bash
example teapot --help
```

## ⭐️ Extending/Improving engine

If you want to play around with the engine and make changes then you want to make sure the unit tests run locally.
To do this you just need to run the following Make command:

```bash
make test
```

This will make sure that the correct `.obj` files are present int `cmd/example/cmd/models` directory.
