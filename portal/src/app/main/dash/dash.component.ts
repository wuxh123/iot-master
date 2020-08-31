import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';

@Component({
  selector: 'app-dash',
  templateUrl: './dash.component.html',
  styleUrls: ['./dash.component.scss']
})
export class DashComponent implements OnInit {

  constructor(private as: ApiService) {
  }

  ngOnInit(): void {

  }

}
