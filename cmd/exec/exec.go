package exec

import (
	"context"
	"io"
	"net/url"

	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// newExecutor returns an Executor. Websocket is preferred, with SPDY as a fallback.
func newExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error) {
	// Try WebSocket first
	exec, err := remotecommand.NewWebSocketExecutor(config, method, url.String())
	if err == nil {
		return exec, nil
	}

	// Fallback to SPDY
	return remotecommand.NewSPDYExecutor(config, method, url)
}

func exec(ctx context.Context, k8sClientSet *kubernetes.Clientset, restClient *restclient.Config, targetPod, targetNamespace string, podCommand []string, fd uintptr, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	request := k8sClientSet.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(targetPod).
		Namespace(targetNamespace).
		SubResource("exec")

	options := &corev1.PodExecOptions{
		Command: podCommand,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	request.VersionedParams(
		options,
		scheme.ParameterCodec,
	)

	executor, err := newExecutor(restClient, "POST", request.URL())
	if err != nil {
		return err
	}

	streamOpts := remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	log.Infof("Connecting to pod %s in namespace %s, use Ctrl+d to exit\n\n", targetPod, targetNamespace)

	// Set the terminal to raw mode
	var previousState *term.State
	previousState, err = term.SetRawTerminal(fd)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the terminal is always restored
	defer term.RestoreTerminal(fd, previousState)
	if err != nil {
		return err
	}

	err = executor.StreamWithContext(ctx, streamOpts)
	if err != nil {
		return err
	}

	return nil
}
