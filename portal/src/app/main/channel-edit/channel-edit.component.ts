import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-channel-edit',
  templateUrl: './channel-edit.component.html',
  styleUrls: ['./channel-edit.component.scss']
})
export class ChannelEditComponent implements OnInit {

  @Input() channel = {};

  constructor() { }

  ngOnInit(): void {
  }

}
