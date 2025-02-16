package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/zpages"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const portEnvKey = "PORT"

func initTracer(zsp *zpages.SpanProcessor) (*sdktrace.TracerProvider, error) {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.New()
	if err != nil {
		return nil, err
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("Proxy"))),
		sdktrace.WithSpanProcessor(zsp),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func main() {
	zsp := zpages.NewSpanProcessor()
	tp, err := initTracer(zsp)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	interceptor := NewInterceptor()

	port := os.Getenv(portEnvKey)
	if port == "" {
		log.Fatalf("%s env var must be set", portEnvKey)
	}

	// Start zPages on a separate port
	go func() {
		http.Handle("/debug/zpages/tracez", zpages.NewTracezHandler(zsp))

		fmt.Println("View tracez at http://localhost:7777/debug/zpages/tracez, public API at http://localhost:7777/hello")
		err = http.ListenAndServe(":7777", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("Proxy server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), interceptor))
}
