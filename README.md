<h1 align="center"> Raytracer </h1>

<p align="center">
Implementation of the ray-tracing in one weekend <a href="https://raytracing.github.io/">book series</a> by Peter Shirley.
<br/>
<br/>
<a href="https://magefile.org">
<img alt="Built with Mage" src="https://img.shields.io/static/v1?label=BUILT%20WITH&message=MAGE&colorA=363a4f&colorB=cba6f7&style=for-the-badge">
</a>
<a href="https://img.shields.io/github/actions/workflow/status/lucasmelin/raytracer/test.yml?colorA=363a4f&colorB=a6da95&label=TESTS&style=for-the-badge">
<img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/lucasmelin/raytracer/test.yml?colorA=363a4f&colorB=a6da95&label=TESTS&style=for-the-badge">
</a>
</p>

## Ray Tracing in One Weekend

[_Ray Tracing in One Weekend_](https://raytracing.github.io/books/RayTracingInOneWeekend.html)

### Final render

![Final render](./outputs/final/oneweekend.png)

## Ray Tracing: The Next Week

[_Ray Tracing: The Next Week_](https://raytracing.github.io/books/RayTracingTheNextWeek.html)

### Final render

![Final render](./outputs/final/thenextweek.png)

## Usage instructions

Install [mage](https://magefile.org/) with Homebrew using `brew install mage`.

### Targets

- `build` - Runs `go mod download` and then builds the `raytracer` binary.
- `clean` - Removes the generated PNG image from disk.
- `install:deps` - Installs all system and Go dependencies.
- `run` - Runs the `raytracer` binary, building it first if necessary.
- `test` - Runs the unit tests.
- `view` - Displays the generated image, generating it first if necessary.
