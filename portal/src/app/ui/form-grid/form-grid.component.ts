import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-form-grid',
  templateUrl: './form-grid.component.html',
  styleUrls: ['./form-grid.component.scss']
})
export class FormGridComponent implements OnInit {
  @Input() label = '';

  constructor() {
  }

  ngOnInit(): void {
  }

}
