import { Service, signal } from '@angular/core';

type ToastKind = 'error' | 'success';

interface LoadingState {
  message: string;
}

interface ToastState {
  kind: ToastKind;
  message: string;
}

interface LoadingRef {
  dismiss(): Promise<void>;
}

@Service()
export class MessagesService {
  readonly loading = signal<LoadingState | null>(null);
  readonly toast = signal<ToastState | null>(null);

  #toastTimeout: ReturnType<typeof setTimeout> | null = null;

  async showLoading(message = 'Working'): Promise<LoadingRef> {
    this.loading.set({ message });
    return {
      dismiss: async () => {
        this.loading.set(null);
      },
    };
  }

  async showErrorToast(message = 'Unexpected error occurred'): Promise<void> {
    this.showToast('error', message);
  }

  async showSuccessToast(message = 'Success'): Promise<void> {
    this.showToast('success', message);
  }

  private showToast(kind: ToastKind, message: string): void {
    if (this.#toastTimeout) {
      clearTimeout(this.#toastTimeout);
    }

    this.toast.set({ kind, message });
    this.#toastTimeout = setTimeout(() => {
      this.toast.set(null);
      this.#toastTimeout = null;
    }, 4000);
  }
}
