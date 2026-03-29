package exec

import (
	"context"
	"io"

	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func exec(ctx context.Context, k8sClientSet *kubernetes.Clientset, restClient *restclient.Config, targetPod, targetNamespace string, podCommand []string, fd uintptr, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	request := k8sClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(targetPod).Namespace(targetNamespace).SubResource("exec")
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

	executor, err := remotecommand.NewSPDYExecutor(restClient, "POST", request.URL())
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
