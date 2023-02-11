import MessageDialog from "../components/dialogs/MessageDialog.vue"

export const showMessageDialog = (modal, message) => {
  modal.show(
    MessageDialog,
    { message }
  );
}