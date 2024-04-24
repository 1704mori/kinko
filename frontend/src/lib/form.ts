export function submit<T = any>(handler: (data: T) => void): (event: SubmitEvent) => void {
	return function factory(event) {
		// @ts-ignore
		const form = event.target!.closest('form') as HTMLFormElement;

		let valid = form.checkValidity();

		const inputs = form.querySelectorAll('input, textarea, .select');

		const payload: Record<string, any> = {};

		for (const input of inputs) {
			// parse 'Select' values
			if (!['INPUT', 'TEXTAREA'].includes(input.nodeName)) {
				let value = input.getAttribute('data-value');
				if (value) {
					if (value.startsWith('[')) {
						value = JSON.parse(value);
					}

					if (!input.getAttribute('name')) {
						continue;
					}

					if (input.getAttribute('required') && !value) {
						valid = false;
					}

					payload[input.getAttribute('name')!] = value;
				} else {
					valid = false;
					validateForm(input as HTMLElement, true);
				}
			} else {
				if (import.meta.env.MODE != 'test') {
					validateForm(input as HTMLElement);
				}
			}
		}

		if (valid) {
			const data = new FormData(form);

			for (const [key, value] of data.entries()) {
				payload[key] = value;
			}

			handler(payload as T);
		} else {
			event.preventDefault();
			if (import.meta.env.MODE != 'test') {
				// @ts-ignore
				[...form].find((element) => element.validity.valid == false)?.focus();
			}
			form.checkValidity();
		}
	};
}

function validateForm(node: HTMLElement, isSelectInvalid?: boolean): { destroy: () => void } {
	// @ts-ignore
	let output = node.tagName === "BUTTON" && node.classList.contains("select") ?  node.parentNode?.previousElementSibling : node.parentNode;

	const errorSpan = getErrorSpan(output!);
	console.log(errorSpan);

	node.addEventListener('invalid', onInvalid);
	node.addEventListener('input', onInput);
	if (isSelectInvalid) {
		onInvalid();
	}

	function onInvalid() {
		// @ts-ignore
		errorSpan.textContent = node.tagName === 'BUTTON' && node.classList.contains('select') ? 'Selecione uma opção' : node.validationMessage;
	}

	function onInput() {
		// @ts-ignore
		if (node.validationMessage == '') errorSpan.textContent = '';
	}

	return {
		destroy() {
			node.removeEventListener('invalid', onInvalid);
			node.removeEventListener('input', onInput);
		}
	};
}

function getErrorSpan(node: HTMLElement) {
	// @ts-ignore
	if (node.firstChild?.firstElementChild?.nodeName === 'SPAN') {
		// @ts-ignore
		return node.firstChild.firstElementChild;
	}

	console.log(node);
	// @ts-ignore
	if (node.lastChild.nodeName === 'SPAN') {
		return node.lastChild;
	}

	// @ts-ignore
	if (node.parentNode.previousElementSibling.lastChild.nodeName === 'SPAN') {
		// @ts-ignore
		return node.parentNode.previousElementSibling.lastChild;
	}
}

function getCustomErrorMessage(node: HTMLInputElement, t: any) {
	if (node.validity.valueMissing) {
		return t.get('auth.errors.valueMissing');
	} else if (node.validity.typeMismatch) {
		return t.get('auth.errors.typeMismatch');
	} else if (node.validity.tooShort) {
		return t.get('auth.errors.tooShort', { minLength: node.minLength });
	} else if (node.validity.tooLong) {
		return t.get('auth.errors.tooLong', { maxLength: node.maxLength });
	} else if (node.validity.patternMismatch) {
		return t.get('auth.errors.patternMismatch');
	} else if (node.validity.rangeOverflow) {
		return t.get('auth.errors.rangeOverflow', { max: node.max });
	} else if (node.validity.rangeUnderflow) {
		return t.get('auth.errors.rangeUnderflow', { min: node.min });
	} else if (node.validity.stepMismatch) {
		return t.get('auth.errors.stepMismatch');
	} else if (node.validity.badInput) {
		return t.get('auth.errors.badInput');
	}

	return '';
}