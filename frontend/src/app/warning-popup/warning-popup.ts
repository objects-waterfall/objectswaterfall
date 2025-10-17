import { Component, input, output } from '@angular/core';

@Component({
  selector: 'app-warning-popup',
  standalone: true,
  templateUrl: './warning-popup.html',
  styleUrls: ['./warning-popup.css', '../../assets/styles/settings-controls.css']
})
export class WarningPopup {
  showPopup = input<boolean>(false);
  title = input<string>("")
  message = input<string>("")
  confirm = output<boolean>()

  confirmAction() {
    this.confirm.emit(true)
  }

  closePopup(){
    this.confirm.emit(false)
  }
}

