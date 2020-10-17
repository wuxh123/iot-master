import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../../api.service';
import {ActivatedRoute} from '@angular/router';
import {TabRef} from '../../tabs/tabs.component';

@Component({
  selector: 'app-model-adapter-edit',
  templateUrl: './model-adapter-edit.component.html',
  styleUrls: ['./model-adapter-edit.component.scss']
})
export class ModelAdapterEditComponent implements OnInit {
  target = 'adapter';
  id = 0;

  data: any = {};

  constructor(private as: ApiService, private routeInfo: ActivatedRoute, private tab: TabRef) {
    tab.name = '链路创建';
  }

  ngOnInit(): void {
    this.id = this.routeInfo.snapshot.params.id || 0;
    if (this.id > 0) {
      this.as.get(this.target + '/' + this.id).subscribe(res => {
        this.data = res.data;
      });
    }
  }

  submit(): void {
    if (this.data.id) {
      this.as.put(this.target + '/' + this.data.id, this.data).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.tab.Close();
      });
    } else {
      this.as.post(this.target, this.data).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.tab.Close();
      });
    }
  }
}
