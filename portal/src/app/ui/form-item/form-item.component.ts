import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-form-item',
  templateUrl: './form-item.component.html',
  styleUrls: ['./form-item.component.scss']
})
export class FormItemComponent implements OnInit {
  @Input() label = '';

  constructor() {
  }

  ngOnInit(): void {
  }

}
